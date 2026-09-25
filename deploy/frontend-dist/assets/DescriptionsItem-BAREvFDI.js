import{Ci as e,Cr as t,Di as n,Ei as r,Fr as i,Ir as a,Oi as o,Si as s,Ti as c,ca as l,di as u,ki as d,la as f,sr as p,ta as m,ti as h,wi as g,wr as _}from"./router-CZSFYymG.js";import{a as v}from"./index-C1zbMUbF.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=s([e(`descriptions`,{fontSize:`var(--n-font-size)`},[e(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),e(`descriptions-table-wrapper`,[e(`descriptions-table`,[e(`descriptions-table-row`,[e(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),e(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),r(`bordered`,[e(`descriptions-table-wrapper`,[e(`descriptions-table`,[e(`descriptions-table-row`,[s(`&:last-child`,[e(`descriptions-table-content`,{paddingBottom:0})])])])])]),c(`left-label-placement`,[e(`descriptions-table-content`,[s(`> *`,{verticalAlign:`top`})])]),c(`left-label-align`,[s(`th`,{textAlign:`left`})]),c(`center-label-align`,[s(`th`,{textAlign:`center`})]),c(`right-label-align`,[s(`th`,{textAlign:`right`})]),c(`bordered`,[e(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[e(`descriptions-table`,[e(`descriptions-table-row`,[s(`&:not(:last-child)`,[e(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),e(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),e(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[s(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),e(`descriptions-table-content`,[s(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),e(`descriptions-header`,`
 font-weight: var(--n-th-font-weight);
 font-size: 18px;
 transition: color .3s var(--n-bezier);
 line-height: var(--n-line-height);
 margin-bottom: 16px;
 color: var(--n-title-text-color);
 `),e(`descriptions-table-wrapper`,`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[e(`descriptions-table`,`
 width: 100%;
 border-collapse: separate;
 border-spacing: 0;
 box-sizing: border-box;
 `,[e(`descriptions-table-row`,`
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[e(`descriptions-table-header`,`
 font-weight: var(--n-th-font-weight);
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-th-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),e(`descriptions-table-content`,`
 vertical-align: top;
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-td-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[g(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),g(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),e(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),o(e(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),d(e(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=l({name:`Descriptions`,props:Object.assign(Object.assign({},p.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(e){let{mergedClsPrefixRef:r,inlineThemeDisabled:i,mergedComponentPropsRef:a}=_(e),o=m(()=>e.size||a?.value?.Descriptions?.size||`medium`),s=p(`Descriptions`,`-descriptions`,b,v,e,r),c=m(()=>{let{bordered:t}=e,r=o.value,{common:{cubicBezierEaseInOut:i},self:{titleTextColor:a,thColor:c,thColorModal:l,thColorPopover:u,thTextColor:d,thFontWeight:f,tdTextColor:p,tdColor:m,tdColorModal:h,tdColorPopover:g,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[n(`fontSize`,r)]:S,[n(t?`thPaddingBordered`:`thPadding`,r)]:C,[n(t?`tdPaddingBordered`:`tdPadding`,r)]:w}}=s.value;return{"--n-title-text-color":a,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":i,"--n-th-font-weight":f,"--n-line-height":x,"--n-th-text-color":d,"--n-td-text-color":p,"--n-th-color":c,"--n-th-color-modal":l,"--n-th-color-popover":u,"--n-td-color":m,"--n-td-color-modal":h,"--n-td-color-popover":g,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),l=i?t(`descriptions`,m(()=>{let t=``,{bordered:n}=e;return n&&(t+=`a`),t+=o.value[0],t}),c,e):void 0;return{mergedClsPrefix:r,cssVars:i?void 0:c,themeClass:l?.themeClass,onRender:l?.onRender,compitableColumn:h(e,[`columns`,`column`]),inlineThemeDisabled:i,mergedSize:o}},render(){let e=this.$slots.default,t=e?a(e()):[];t.length;let{contentClass:n,labelClass:r,compitableColumn:o,labelPlacement:s,labelAlign:c,mergedSize:l,bordered:d,title:p,cssVars:m,mergedClsPrefix:h,separator:g,onRender:_}=this;_?.();let v=t.filter(e=>S(e)),b=v.reduce((e,t,i)=>{let a=t.props||{},c=v.length-1===i,l=[`label`in a?a.label:y(t,`label`)],u=[y(t)],p=a.span||1,m=e.span;e.span+=p;let _=a.labelStyle||a[`label-style`]||this.labelStyle,b=a.contentStyle||a[`content-style`]||this.contentStyle;if(s===`left`)d?e.row.push(f(`th`,{class:[`${h}-descriptions-table-header`,r],colspan:1,style:_},l),f(`td`,{class:[`${h}-descriptions-table-content`,n],colspan:c?(o-m)*2+1:p*2-1,style:b},u)):e.row.push(f(`td`,{class:`${h}-descriptions-table-content`,colspan:c?(o-m)*2:p*2},f(`span`,{class:[`${h}-descriptions-table-content__label`,r],style:_},[...l,g&&f(`span`,{class:`${h}-descriptions-separator`},g)]),f(`span`,{class:[`${h}-descriptions-table-content__content`,n],style:b},u)));else{let t=c?(o-m)*2:p*2;e.row.push(f(`th`,{class:[`${h}-descriptions-table-header`,r],colspan:t,style:_},l)),e.secondRow.push(f(`td`,{class:[`${h}-descriptions-table-content`,n],colspan:t,style:b},u))}return(e.span>=o||c)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),s!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>f(`tr`,{class:`${h}-descriptions-table-row`},e));return f(`div`,{style:m,class:[`${h}-descriptions`,this.themeClass,`${h}-descriptions--${s}-label-placement`,`${h}-descriptions--${c}-label-align`,`${h}-descriptions--${l}-size`,d&&`${h}-descriptions--bordered`]},p||this.$slots.header?f(`div`,{class:`${h}-descriptions-header`},p||i(this,`header`)):null,f(`div`,{class:`${h}-descriptions-table-wrapper`},f(`table`,{class:`${h}-descriptions-table`},f(`tbody`,null,s===`top`&&f(`tr`,{class:`${h}-descriptions-table-row`,style:{visibility:`collapse`}},u(o*2,f(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=l({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};