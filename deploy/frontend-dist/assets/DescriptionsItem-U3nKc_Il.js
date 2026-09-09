import{Ar as e,Ci as t,Si as n,Yi as r,Yr as i,_i as a,_r as o,ai as s,bi as c,gi as l,kr as u,na as d,ra as f,tr as p,vi as m,vr as h,xi as g,yi as _}from"./router-DqBWDa9m.js";import{a as v}from"./index-0gynlAz9.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=l([a(`descriptions`,{fontSize:`var(--n-font-size)`},[a(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),a(`descriptions-table-wrapper`,[a(`descriptions-table`,[a(`descriptions-table-row`,[a(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),a(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),c(`bordered`,[a(`descriptions-table-wrapper`,[a(`descriptions-table`,[a(`descriptions-table-row`,[l(`&:last-child`,[a(`descriptions-table-content`,{paddingBottom:0})])])])])]),_(`left-label-placement`,[a(`descriptions-table-content`,[l(`> *`,{verticalAlign:`top`})])]),_(`left-label-align`,[l(`th`,{textAlign:`left`})]),_(`center-label-align`,[l(`th`,{textAlign:`center`})]),_(`right-label-align`,[l(`th`,{textAlign:`right`})]),_(`bordered`,[a(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[a(`descriptions-table`,[a(`descriptions-table-row`,[l(`&:not(:last-child)`,[a(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),a(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),a(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[l(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),a(`descriptions-table-content`,[l(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),a(`descriptions-header`,`
 font-weight: var(--n-th-font-weight);
 font-size: 18px;
 transition: color .3s var(--n-bezier);
 line-height: var(--n-line-height);
 margin-bottom: 16px;
 color: var(--n-title-text-color);
 `),a(`descriptions-table-wrapper`,`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[a(`descriptions-table`,`
 width: 100%;
 border-collapse: separate;
 border-spacing: 0;
 box-sizing: border-box;
 `,[a(`descriptions-table-row`,`
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[a(`descriptions-table-header`,`
 font-weight: var(--n-th-font-weight);
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-th-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),a(`descriptions-table-content`,`
 vertical-align: top;
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-td-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[m(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),m(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),a(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),n(a(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),t(a(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=d({name:`Descriptions`,props:Object.assign(Object.assign({},p.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n,mergedComponentPropsRef:a}=h(e),s=r(()=>e.size||a?.value?.Descriptions?.size||`medium`),c=p(`Descriptions`,`-descriptions`,b,v,e,t),l=r(()=>{let{bordered:t}=e,n=s.value,{common:{cubicBezierEaseInOut:r},self:{titleTextColor:i,thColor:a,thColorModal:o,thColorPopover:l,thTextColor:u,thFontWeight:d,tdTextColor:f,tdColor:p,tdColorModal:m,tdColorPopover:h,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[g(`fontSize`,n)]:S,[g(t?`thPaddingBordered`:`thPadding`,n)]:C,[g(t?`tdPaddingBordered`:`tdPadding`,n)]:w}}=c.value;return{"--n-title-text-color":i,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":r,"--n-th-font-weight":d,"--n-line-height":x,"--n-th-text-color":u,"--n-td-text-color":f,"--n-th-color":a,"--n-th-color-modal":o,"--n-th-color-popover":l,"--n-td-color":p,"--n-td-color-modal":m,"--n-td-color-popover":h,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),u=n?o(`descriptions`,r(()=>{let t=``,{bordered:n}=e;return n&&(t+=`a`),t+=s.value[0],t}),l,e):void 0;return{mergedClsPrefix:t,cssVars:n?void 0:l,themeClass:u?.themeClass,onRender:u?.onRender,compitableColumn:i(e,[`columns`,`column`]),inlineThemeDisabled:n,mergedSize:s}},render(){let t=this.$slots.default,n=t?e(t()):[];n.length;let{contentClass:r,labelClass:i,compitableColumn:a,labelPlacement:o,labelAlign:c,mergedSize:l,bordered:d,title:p,cssVars:m,mergedClsPrefix:h,separator:g,onRender:_}=this;_?.();let v=n.filter(e=>S(e)),b=v.reduce((e,t,n)=>{let s=t.props||{},c=v.length-1===n,l=[`label`in s?s.label:y(t,`label`)],u=[y(t)],p=s.span||1,m=e.span;e.span+=p;let _=s.labelStyle||s[`label-style`]||this.labelStyle,b=s.contentStyle||s[`content-style`]||this.contentStyle;if(o===`left`)d?e.row.push(f(`th`,{class:[`${h}-descriptions-table-header`,i],colspan:1,style:_},l),f(`td`,{class:[`${h}-descriptions-table-content`,r],colspan:c?(a-m)*2+1:p*2-1,style:b},u)):e.row.push(f(`td`,{class:`${h}-descriptions-table-content`,colspan:c?(a-m)*2:p*2},f(`span`,{class:[`${h}-descriptions-table-content__label`,i],style:_},[...l,g&&f(`span`,{class:`${h}-descriptions-separator`},g)]),f(`span`,{class:[`${h}-descriptions-table-content__content`,r],style:b},u)));else{let t=c?(a-m)*2:p*2;e.row.push(f(`th`,{class:[`${h}-descriptions-table-header`,i],colspan:t,style:_},l)),e.secondRow.push(f(`td`,{class:[`${h}-descriptions-table-content`,r],colspan:t,style:b},u))}return(e.span>=a||c)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),o!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>f(`tr`,{class:`${h}-descriptions-table-row`},e));return f(`div`,{style:m,class:[`${h}-descriptions`,this.themeClass,`${h}-descriptions--${o}-label-placement`,`${h}-descriptions--${c}-label-align`,`${h}-descriptions--${l}-size`,d&&`${h}-descriptions--bordered`]},p||this.$slots.header?f(`div`,{class:`${h}-descriptions-header`},p||u(this,`header`)):null,f(`div`,{class:`${h}-descriptions-table-wrapper`},f(`table`,{class:`${h}-descriptions-table`},f(`tbody`,null,o===`top`&&f(`tr`,{class:`${h}-descriptions-table-row`,style:{visibility:`collapse`}},s(a*2,f(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=d({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};