import{Ci as e,Ei as t,Mr as n,Nr as r,Qi as i,Qr as a,Si as o,Ti as s,aa as c,bi as l,br as u,ci as d,ir as f,oa as p,wi as m,xi as h,xr as g,yi as _}from"./router-Ck3NDp5z.js";import{a as v}from"./index-AfrKiTqj.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=_([l(`descriptions`,{fontSize:`var(--n-font-size)`},[l(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),l(`descriptions-table-wrapper`,[l(`descriptions-table`,[l(`descriptions-table-row`,[l(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),l(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),e(`bordered`,[l(`descriptions-table-wrapper`,[l(`descriptions-table`,[l(`descriptions-table-row`,[_(`&:last-child`,[l(`descriptions-table-content`,{paddingBottom:0})])])])])]),o(`left-label-placement`,[l(`descriptions-table-content`,[_(`> *`,{verticalAlign:`top`})])]),o(`left-label-align`,[_(`th`,{textAlign:`left`})]),o(`center-label-align`,[_(`th`,{textAlign:`center`})]),o(`right-label-align`,[_(`th`,{textAlign:`right`})]),o(`bordered`,[l(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[l(`descriptions-table`,[l(`descriptions-table-row`,[_(`&:not(:last-child)`,[l(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),l(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),l(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[_(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),l(`descriptions-table-content`,[_(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),l(`descriptions-header`,`
 font-weight: var(--n-th-font-weight);
 font-size: 18px;
 transition: color .3s var(--n-bezier);
 line-height: var(--n-line-height);
 margin-bottom: 16px;
 color: var(--n-title-text-color);
 `),l(`descriptions-table-wrapper`,`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[l(`descriptions-table`,`
 width: 100%;
 border-collapse: separate;
 border-spacing: 0;
 box-sizing: border-box;
 `,[l(`descriptions-table-row`,`
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[l(`descriptions-table-header`,`
 font-weight: var(--n-th-font-weight);
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-th-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),l(`descriptions-table-content`,`
 vertical-align: top;
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-td-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[h(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),h(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),l(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),s(l(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),t(l(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=c({name:`Descriptions`,props:Object.assign(Object.assign({},f.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n,mergedComponentPropsRef:r}=g(e),o=i(()=>e.size||r?.value?.Descriptions?.size||`medium`),s=f(`Descriptions`,`-descriptions`,b,v,e,t),c=i(()=>{let{bordered:t}=e,n=o.value,{common:{cubicBezierEaseInOut:r},self:{titleTextColor:i,thColor:a,thColorModal:c,thColorPopover:l,thTextColor:u,thFontWeight:d,tdTextColor:f,tdColor:p,tdColorModal:h,tdColorPopover:g,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[m(`fontSize`,n)]:S,[m(t?`thPaddingBordered`:`thPadding`,n)]:C,[m(t?`tdPaddingBordered`:`tdPadding`,n)]:w}}=s.value;return{"--n-title-text-color":i,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":r,"--n-th-font-weight":d,"--n-line-height":x,"--n-th-text-color":u,"--n-td-text-color":f,"--n-th-color":a,"--n-th-color-modal":c,"--n-th-color-popover":l,"--n-td-color":p,"--n-td-color-modal":h,"--n-td-color-popover":g,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),l=n?u(`descriptions`,i(()=>{let t=``,{bordered:n}=e;return n&&(t+=`a`),t+=o.value[0],t}),c,e):void 0;return{mergedClsPrefix:t,cssVars:n?void 0:c,themeClass:l?.themeClass,onRender:l?.onRender,compitableColumn:a(e,[`columns`,`column`]),inlineThemeDisabled:n,mergedSize:o}},render(){let e=this.$slots.default,t=e?r(e()):[];t.length;let{contentClass:i,labelClass:a,compitableColumn:o,labelPlacement:s,labelAlign:c,mergedSize:l,bordered:u,title:f,cssVars:m,mergedClsPrefix:h,separator:g,onRender:_}=this;_?.();let v=t.filter(e=>S(e)),b=v.reduce((e,t,n)=>{let r=t.props||{},c=v.length-1===n,l=[`label`in r?r.label:y(t,`label`)],d=[y(t)],f=r.span||1,m=e.span;e.span+=f;let _=r.labelStyle||r[`label-style`]||this.labelStyle,b=r.contentStyle||r[`content-style`]||this.contentStyle;if(s===`left`)u?e.row.push(p(`th`,{class:[`${h}-descriptions-table-header`,a],colspan:1,style:_},l),p(`td`,{class:[`${h}-descriptions-table-content`,i],colspan:c?(o-m)*2+1:f*2-1,style:b},d)):e.row.push(p(`td`,{class:`${h}-descriptions-table-content`,colspan:c?(o-m)*2:f*2},p(`span`,{class:[`${h}-descriptions-table-content__label`,a],style:_},[...l,g&&p(`span`,{class:`${h}-descriptions-separator`},g)]),p(`span`,{class:[`${h}-descriptions-table-content__content`,i],style:b},d)));else{let t=c?(o-m)*2:f*2;e.row.push(p(`th`,{class:[`${h}-descriptions-table-header`,a],colspan:t,style:_},l)),e.secondRow.push(p(`td`,{class:[`${h}-descriptions-table-content`,i],colspan:t,style:b},d))}return(e.span>=o||c)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),s!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>p(`tr`,{class:`${h}-descriptions-table-row`},e));return p(`div`,{style:m,class:[`${h}-descriptions`,this.themeClass,`${h}-descriptions--${s}-label-placement`,`${h}-descriptions--${c}-label-align`,`${h}-descriptions--${l}-size`,u&&`${h}-descriptions--bordered`]},f||this.$slots.header?p(`div`,{class:`${h}-descriptions-header`},f||n(this,`header`)):null,p(`div`,{class:`${h}-descriptions-table-wrapper`},p(`table`,{class:`${h}-descriptions-table`},p(`tbody`,null,s===`top`&&p(`tr`,{class:`${h}-descriptions-table-row`,style:{visibility:`collapse`}},d(o*2,p(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=c({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};